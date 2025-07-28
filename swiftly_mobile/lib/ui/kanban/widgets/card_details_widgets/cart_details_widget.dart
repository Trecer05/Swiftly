import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/core/ui/label_item_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/card_details_widgets/status_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/date_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/priority_widget.dart';

import '../../../../domain/kanban/models/card_item.dart';
import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/label_notifier_provider.dart';
import '../../../../providers/user_notifier_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import '../../../home/widgets/content/avatar_widget.dart';
import 'button_widget.dart';
import 'text_field_widget.dart';

class CartDetailsWidget extends ConsumerWidget {
  final CardItem card;

  CartDetailsWidget({super.key, required this.card});

  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  late String _selectedColumnId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    _titleController.text = card.title;
    _descriptionController.text = card.description;
    _selectedColumnId = card.columnId;

    final labels =
        ref
            .watch(labelNotifierProvider)
            .labels
            .where((label) => label.cardId == card.id)
            .toList();

    final users = ref.watch(userNotifierProvider).users;
    final currentUser = users.firstWhere((user) => user.id == card.userId);
    return Center(
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
          child: SingleChildScrollView(
            child: Container(
              constraints: const BoxConstraints(maxWidth: 400),
              width: 400,
              height: 500,
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
              decoration: BoxDecoration(
                color: AppColors.white26,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      SizedBox(
                        width: 300,
                        child: TextFieldWidget(
                          hintText: 'Название карточки',
                          controller: _titleController,
                        ),
                      ),
                      const Spacer(),
                      DateWidget(
                        date: card.createdAt,
                        color: AppColors.twitter,
                      ),
                    ],
                  ),
                  const SizedBox(height: 10),
                  SizedBox(
                    width: 300,
                    child: TextFieldWidget(
                      hintText: 'Описание',
                      controller: _descriptionController,
                    ),
                  ),
                  const SizedBox(height: 10),
                  Row(
                    children: [
                      const Text(
                        'Статус',
                        style: TextStyle(color: AppColors.white),
                      ),
                      const SizedBox(width: 10),
                      StatusWidget(
                        statusId: _selectedColumnId,
                        onChanged: (newId) {
                          _selectedColumnId = newId;
                        },
                      ),
                    ],
                  ),
                  const SizedBox(height: 10),
                  if (labels.isNotEmpty)
                    Column(
                      children: [
                        Row(
                          children:
                              labels
                                  .map(
                                    (l) => Padding(
                                      padding: const EdgeInsets.only(right: 8),
                                      child: LabelItemWidget(labelItem: l),
                                    ),
                                  )
                                  .toList(),
                        ),
                        const SizedBox(height: 10),
                      ],
                    ),
                  if (card.priority != null)
                    Column(
                      children: [
                        PriorityWidget(priority: card.priority!),
                        const SizedBox(height: 10),
                      ],
                    ),
                  // Row(
                  //   crossAxisAlignment: CrossAxisAlignment.start,
                  //   children: [
                  //     AvatarWidget(imageUrl: currentUser.image),
                  //     const SizedBox(width: 5),
                  //     Expanded(
                  //       child: Column(
                  //         crossAxisAlignment: CrossAxisAlignment.start,
                  //         children: [
                  //           Text(
                  //             '${currentUser.name} ${currentUser.lastName ?? ''}',
                  //             style: AppTextStyles.style13,
                  //             overflow: TextOverflow.ellipsis,
                  //             maxLines: 1,
                  //             softWrap: false,
                  //           ),
                  //           if (currentUser.role != null)
                  //             LabelItemWidget(labelItem: currentUser.role!),
                  //         ],
                  //       ),
                  //     ),
                  //   ],
                  // ),
                  const Spacer(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      ButtonWidget(
                        title: 'Удалить',
                        color: AppColors.red,
                        onTap: () => _handleDelete(context, ref, card.id),
                      ),
                      const SizedBox(width: 10),
                      ButtonWidget(
                        title: 'Сохранить',
                        color: AppColors.blue,
                        onTap: () => _handleSave(context, ref, card.id),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _handleDelete(BuildContext context, WidgetRef ref, String id) {
    ref.read(cardNotifierProvider.notifier).removeCard(id);
    Navigator.of(context).pop();
  }

  void _handleSave(BuildContext context, WidgetRef ref, String id) {
    final newTitle = _titleController.text;
    final newDescription = _descriptionController.text;
    final newColumn = _selectedColumnId;
    ref.read(cardNotifierProvider.notifier).updateTitle(id, newTitle);
    ref
        .read(cardNotifierProvider.notifier)
        .updateDescription(id, newDescription);
    ref.read(cardNotifierProvider.notifier).updateColumn(id, newColumn);
    Navigator.of(context).pop();
  }
}
